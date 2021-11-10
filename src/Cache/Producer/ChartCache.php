<?php
/**
 * ChartCache.php
 * @author      djh <dengjinghui@shiyue.com>
 * @date     2021/9/17 11:57
 * PhpStorm
 */


namespace CacheSystem\Cache\Producer;


use CacheSystem\Cache\Builder\AbstractCache;
use CacheSystem\Cache\Builder\RedisCache;
use CacheSystem\Cache\Exception\CacheException;
use CacheSystem\Facade\Query;
use CacheSystem\Query\Producer\QueryInterface;
use Illuminate\Support\Facades\App;

class ChartCache extends SysCache
{

    /**
     * Construct a query operation object
     *
     * @param QueryInterface|null $queryObject
     * @return static
     */
    public static function create(QueryInterface $queryDriver = null)
    {
        if ($queryDriver === null) {
            $self = new static();
            $self->queryDriver = Query::ImpalaQuery();
            return $self;
        }

        return new static($queryDriver);
    }

    /**
     * Construct a cache operation object
     *
     * @param CacheAbstract|null $cacheObject
     * @return $this
     */
    public function build(AbstractCache $cacheObject = null)
    {
        if ($cacheObject === null)
            $this->cacheDriver = RedisCache::make();
        else
            $this->cacheDriver = $cacheObject;

        return $this;
    }

    public function get($sql)
    {

        $this->debug('缓存系统启动.');
        if ($this->cacheDriver === null)
            throw new CacheException("Cache driver is not initialized.");

        if (!$sql)
            throw new CacheException("Missing sql parameter.");

        if (is_callable($sql))
            $sql = call_user_func($sql);

        try {

            $this->setSql($sql);

            $this->cacheDriver->setSql($sql);

            $this->cacheDriver->setCacheKey($this->cacheDriver::parseCacheKey($sql));
            $this->debug("初始化缓存KEY成功.");

            if ($this->executeDelete) {

                $this->debug("开始清除缓存.");

                $this->cacheDriver->clear($sql);

                $this->debug("清除缓存成功.");
            }

            $cacheResult = [];

            $cacheAttr = $this->cacheDriver->resetConfig($this);

            if ($cacheAttr->isCache) {

                $this->debug('开始获取缓存.');

                $cacheResult = $this->timeConsumingCalculation(function ($startTime) use ($sql) {
                    $result = $this->cacheDriver->get($sql);
                    $this->cacheDriver->setUseTime(microtime(true) - $startTime);
                    return $result;
                });

                $this->debug(function () {
                    if ($notHitSql = (array)$this->cacheDriver->getSql()) {
                        $notHitSqlKey = [];
                        array_walk($notHitSql, function ($sql, $key) use (&$notHitSqlKey) {
                            $notHitSqlKey[$key] = $this->cacheDriver->generateKey($sql);
                        });
                    }
                    if (!empty($notHitSqlKey))
                        return sprintf("获取缓存结束,未获取key值:%s.", json_encode($notHitSqlKey));
                    return sprintf("获取缓存结束.是否获取到数据:%s.", empty($this->cacheDriver->getSql()) ? '是' : '否');
                });
            }
        } catch (\Throwable $throwable) {

            $this->debug("发生异常：{$throwable->__toString()}");

            if (App::environment() != 'local')
                \Log::error("发生异常：{$throwable->__toString()}");

            $cacheResult = $this->cacheDriver->responseFail($sql);

            $cacheAttr = $this;
        }

        $queryResult = [];

        if (!empty($this->cacheDriver->getSql())) {

            $this->debug('开始查询数据库.');

            $queryResult = $this->timeConsumingCalculation(function ($startTime) use ($sql) {
                $result = $this->queryDriver->get($this->cacheDriver->getSql());
                $this->setUseTime(microtime(true) - $startTime);
                return $result;
            });

            $this->debug('查询数据库结束.');

            if ($cacheAttr->isCache && $this->validate($queryResult)) {

                try {
                    $this->debug("开始写入缓存数据.");
                    //开始缓存
//                    if ($this->cacheResultSet && is_callable($this->cacheResultSet)) {
//                        $result = call_user_func($this->cacheResultSet, $result);
//                    }
                    $this->cacheDriver->save($cacheAttr, $queryResult);

                } catch (\Throwable $throwable) {

                    $this->debug("写入缓存时发生异常：{$throwable->__toString()}");

                    if (App::environment() != 'local')
                        \Log::error("写入缓存时发生异常：{$throwable->__toString()}");
                }

                $this->debug("写入缓存数据成功！.");
            }
        }

        $this->debug("通知观察者.");
        $this->notify($this, $cacheAttr);

        return array_merge($cacheResult, $queryResult);

    }

    public function save($cache, $result)
    {
        // TODO: Implement save() method.
    }


    public function clear($sql)
    {
        // TODO: Implement clear() method.
    }


}
