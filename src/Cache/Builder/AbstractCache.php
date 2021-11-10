<?php
/**
 * abstractCache.php
 * @author      djh <dengjinghui@shiyue.com>
 * @date     2021/9/17 14:15
 * PhpStorm
 */


namespace CacheSystem\Cache\Builder;


use App\Models\CacheSystem;
use App\Models\CacheSystemConfig;
use CacheSystem\Cache\Producer\SysCache;

/**
 * Class AbstractCache
 * @package CacheSystem\Cache\Builder
 */
abstract class AbstractCache extends SysCache
{
    /**
     * Whether to hit the cache
     *
     * @var bool
     */
    protected $cacheHit = false;

    /**
     * The size of the result set
     *
     * @var string|array
     */
    protected $dataSize;

    /**
     * Cache key
     *
     * @var string|array
     */
    protected $cacheKey;

    /**
     * AbstractCache constructor.
     */
    private function __construct()
    {
        $this->init();
    }

    protected function init()
    {

    }

    /**
     * @return static
     */
    public static function make()
    {
        return new static();
    }

    /**
     * @return bool
     */
    public function isCacheHit(): bool
    {
        return $this->cacheHit;
    }

    /**
     * @param bool $cacheHit
     */
    public function setCacheHit(bool $cacheHit): void
    {
        $this->cacheHit = $cacheHit;
    }


    /**
     * @return array|string
     */
    public function getCacheKey()
    {
        return $this->cacheKey;
    }

    /**
     * @param array|string $cacheKey
     */
    public function setCacheKey($cacheKey): void
    {
        $this->cacheKey = $cacheKey;
    }

    /**
     * @return mixed
     */
    public function getDataSize()
    {
        return $this->dataSize;
    }

    /**
     * @param int $dataSize
     */
    public function setDataSize(int $dataSize): void
    {
        $this->dataSize = $dataSize;
    }


    /**
     * @param null $sql
     * @return array
     */
    public function responseFail($sql = null)
    {
        $this->cacheHit = false;
        $this->sql = $sql;
        return [];
    }

    /**
     * @param $result
     * @return mixed
     */
    public function responseSuccess($result)
    {
        $this->cacheHit = true;
        $this->sql = null;
        return is_array($result) ? $result : \GuzzleHttp\json_decode($result, true);
    }

    public static function parseCacheKey($sql)
    {
        if (is_string($sql))
            return self::removeSpaces($sql);

        if (is_array($sql))
            return array_map(function ($sql) {
                return self::removeSpaces($sql);
            }, $sql);

        return null;
    }

    private static function removeSpaces(string $sql)
    {
        return md5(preg_replace('/\s+/', '', $sql));
    }

    /**
     * @param $sql
     * @return array|string|string[]|null
     */
    public function generateKey($sql)
    {

        if (is_string($sql))
            return md5($sql);

        if (is_array($sql))
            return array_map(function ($sql) {
                return md5($sql);
            }, $sql);

        return null;
    }

    protected function initialization($sql)
    {
        $this->sql = $sql;
    }

    /**
     * @param SysCache $cache
     * @return SysCache|\Illuminate\Database\Eloquent\Builder|\Illuminate\Database\Eloquent\Model|object|null
     */
    public function resetConfig(SysCache $cache)
    {
        try {
            if ($cache->dbConfig)
                if ($dbConfig = CacheSystemConfig::query()->where('cacheKey', $cache->cacheDriver->cacheKey)->first())
                    return $dbConfig;

            return $cache;
        } catch (\Exception $exception) {

            return $cache;
        }

    }

    /**
     * Calculate the data size in kb
     *
     * @param string $data
     * @return float|int
     */
    protected function calculateTheSize(string $data)
    {
        return strlen($data) * 8 / 1024;
    }

}
