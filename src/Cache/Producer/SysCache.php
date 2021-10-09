<?php
/**
 * SysCache.php
 * @author      djh <dengjinghui@shiyue.com>
 * @date     2021/9/17 11:46
 * PhpStorm
 */

namespace CacheSystem\Cache\Producer;

use App\Components\Helper\Debug;
use CacheSystem\Cache\Builder\AbstractCache;
use CacheSystem\Cache\Exception\CacheException;
use CacheSystem\Cache\Observer\CacheObserver;
use CacheSystem\Cache\Observer\SaveHistorySqlImpl;
use CacheSystem\Cache\Validate\Config;
use CacheSystem\Cache\Validate\ValidationException;
use CacheSystem\Cache\Validate\ValidationInterface;
use CacheSystem\Facade\Query;
use \CacheSystem\Query\Producer\QueryInterface;
use Illuminate\Support\Facades\App;
use Illuminate\Support\Facades\Request;

abstract class SysCache
{
    use Debug;

    /**
     * API mode, will output to the log file
     */
    const DEBUG_MODE_API = 1;

    /**
     * WEB mode, will be output directly on the page
     */
    const DEBUG_MODE_WEB = 2;

    /**
     * Cache status, whether to cache
     *
     * @var bool
     */
    protected $isCache = true;

    protected $sql;

    /**
     * Log output mode
     *
     * @var int
     */
    protected $debugMode = self::DEBUG_MODE_WEB;

    /**
     * Program start time
     *
     * @var float|string
     */
    protected $useTime = 0;

    /**
     * Cache valid time 5 minutes by default(s)
     *
     * @var int
     */
    protected $expTime = 300;

    /**
     * Query-driven, used to interact with the database
     *
     * @var QueryInterface
     */
    protected $queryDriver;

    /**
     * Cache-driven
     *
     * @var AbstractCache
     */
    protected $cacheDriver;

    /**
     * Observer collection
     *
     * @var array
     */
    protected $observerMap = [
        SaveHistorySqlImpl::class
    ];

    /**
     * 开启手动删除缓存状态
     * @var bool
     */
    protected $executeDelete = false;

    /**
     * 可缓存的最大数据行数，默认一万行
     *
     * @var int
     */
    protected $maxCacheDataRows = 10000;

    /**
     * This type of object attributes preferentially use database configuration
     *
     * @var bool
     */
    protected $dbConfig = true;

    /**
     * 缓存结果集操作，如果需要缓存结果，请传入操作结果的逻辑
     *
     * @var callable|null
     */
    private $cacheResultSet = null;


    /**
     * @var int
     */
    private $step = 1;

    /**
     * ChartCache constructor.
     */
    public function __construct()
    {
        if (Request::get('del_cache', false))
            $this->executeDelete = true;
    }

    /**
     * @param array|string $sql
     * @return mixed
     */
    abstract public function get(array $sql);

    /**
     * @param $cache
     * @param $result
     * @return mixed
     */
    abstract public function save($cache, $result);

    /**
     * @return mixed
     */
    abstract public function clear($sql);

    /**
     * @return bool
     */
    public function isDebug(): bool
    {
        return $this->debug;
    }


    /**
     * @param bool $debug
     * @param int $mode
     * @return $this
     */
    public function setDebug(bool $debug, $mode = self::DEBUG_MODE_WEB): self
    {
        $this->debug = $debug;

        $this->debugMode = $mode;

        return $this;
    }

    /**
     * @return bool
     */
    public function isCache(): bool
    {
        return $this->isCache;
    }

    /**
     * @param bool $isCache
     */
    public function setIsCache(bool $isCache): self
    {
        $this->isCache = $isCache;
        return $this;
    }

    /**
     * @return QueryInterface
     */
    public function getQueryDriver(): QueryInterface
    {
        return $this->queryDriver;
    }

    /**
     * @param QueryInterface $queryDriver
     */
    public function setQueryDriver(QueryInterface $queryDriver): void
    {
        $this->queryDriver = $queryDriver;
    }

    /**
     * @return mixed
     */
    public function getCacheDriver()
    {
        return $this->cacheDriver;
    }

    /**
     * @param mixed $cacheDriver
     */
    public function setCacheDriver($cacheDriver): void
    {
        $this->cacheDriver = $cacheDriver;
    }

    /**
     * @return array
     */
    public function getObserverMap(): array
    {
        return $this->observerMap;
    }

    /**
     * @param array $observerMap
     */
    public function setObserverMap(array $observerMap): void
    {
        $this->observerMap = $observerMap;
    }

    /**
     * @return int
     */
    public function getMaxCacheDataRows(): int
    {
        return $this->maxCacheDataRows;
    }

    /**
     * @param int $maxCacheDataRows
     */
    public function setMaxCacheDataRows(int $maxCacheDataRows): void
    {
        $this->maxCacheDataRows = $maxCacheDataRows;
    }

    /**
     * @return callable|null
     */
    public function getCacheResultSet(): ?callable
    {
        return $this->cacheResultSet;
    }

    /**
     * @param callable|null $cacheResultSet
     */
    public function setCacheResultSet(?callable $cacheResultSet): void
    {
        $this->cacheResultSet = $cacheResultSet;
    }

    /**
     * @return mixed
     */
    public function getSql()
    {
        return $this->sql;
    }

    /**
     * @param mixed $sql
     */
    public function setSql($sql): void
    {
        $this->sql = $sql;
    }

    /**
     * @return int
     */
    public function getExpTime(): int
    {
        return $this->expTime;
    }

    /**
     * @param int $expTime
     * @return $this
     */
    public function setExpTime(int $expTime): self
    {
        $this->expTime = $expTime;
        return $this;
    }

    /**
     * @return mixed
     */
    public function getUseTime()
    {
        return $this->useTime;
    }

    /**
     * @param float|string $useTime
     */
    public function setUseTime($useTime): void
    {
        $this->useTime = number_format($useTime, 6);
    }


    /**
     * Observer after cache write result
     *
     * @param mixed ...$params
     */
    protected function notify(...$params)
    {
        array_map(function ($observer) use ($params) {
            (new $observer)->onChange($params);
        }, $this->observerMap);
    }

    /**
     * Validator at cache write
     *
     * @param $result
     * @return bool
     * @throws ValidationException
     */
    protected function validate(&$result): bool
    {
        try {
            array_map(function ($verifier) use (&$result) {
                /** @var ValidationInterface $verifier */
                $verifier = new $verifier;
                $verifier->rule($this, $result);
            }, Config::verifier());
            return true;
        } catch (ValidationException $exception) {

            $this->debug("校验缓存结果时发生异常：{$exception->__toString()}");

            if (App::environment() != 'local')
                \Log::error("校验缓存结果时发生异常：{$exception->__toString()}");

            return false;
        }
    }


    /**
     * @param $func
     */
    protected function debug($func)
    {
        if ($this->debug) {

            if (is_callable($func))
                $func = call_user_func($func);

            if (is_string($func)) {

                if ($this->debugMode == self::DEBUG_MODE_WEB) {

                    print_r(sprintf("<pre><span style='color: red'>[%s]INFO:</span>step%d.%s</pre>\r\n", date("Y-m-d H:i:s.u")
                        , $this->step, $func));

                } else

                    \Log::info(sprintf('step%d.%s', $this->step, $func));

                $this->step++;
            }
        }

    }

}
