<?php
/**
 * RedisCache.php
 * @author      djh <dengjinghui@shiyue.com>
 * @date     2021/9/17 14:21
 * PhpStorm
 */


namespace CacheSystem\Cache\Builder;


use App\Models\CacheSystemConfig;
use CacheSystem\Cache\Producer\SysCache;
use CacheSystem\Cache\Singleton\RedisSingleton;
use Illuminate\Support\Arr;

class RedisCache extends AbstractCache
{
    /**
     * @var \Predis\ClientInterface|null
     */
    protected $db;


    protected function init()
    {
        $this->db = RedisSingleton::instance();
    }

    /**
     * @inheritDoc
     */
    public function get(array $sql)
    {
        $result = [];

        array_walk($this->cacheKey, function ($cacheKey, $key) use (&$result) {
            if ($this->db->exists($cacheKey)) {
                $result[$key] = $this->db->get($cacheKey) ?: '[]';
                unset($this->sql[$key]);
                $this->dataSize[$key] = $this->calculateTheSize($result[$key]);
                $result[$key] = \GuzzleHttp\json_decode($result[$key], true);
            }
        });
        return $result;
    }

    /**
     * @param SysCache|CacheSystemConfig $cache
     * @param $result
     * @return mixed|void
     */
    public function save($cache, $result)
    {
        $result = (array)$result;
        //Single sql cache
        if (is_string($this->cacheKey)) {
            $this->dataSize = $this->calculateTheSize(\GuzzleHttp\json_encode($result));
            $this->db->set($this->cacheKey, \GuzzleHttp\json_encode($result), "EX", $cache->expTime);
            return;
        }

        //Multiple sql caches
        if (is_array($this->sql)) {
            array_walk($this->sql, function ($cacheKey, $key) use ($result, $cache) {
                $this->dataSize[$key] = $this->calculateTheSize(\GuzzleHttp\json_encode($result[$key]));
                $this->db->set($this->generateKey($cacheKey), \GuzzleHttp\json_encode($result[$key]), "EX",
                    $cache->expTime);
            });

        }
    }


    /**
     * @inheritDoc
     */
    public function clear($sql)
    {
        if (empty($sql)) return;

        $this->cacheKey = $this->generateKey($sql);

        if (is_array($this->cacheKey))
            foreach ($this->cacheKey as $cacheKey) {
                $this->db->del($cacheKey);
            }

        if (is_string($this->cacheKey))
            $this->db->del($this->cacheKey);

        return;
    }

    /**
     *
     */
    public function deleteDb()
    {
        $this->db = null;
    }

}
