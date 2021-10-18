<?php
/**
 * RedisSingleton.php
 * @author      djh <dengjinghui@shiyue.com>
 * @date     2021/9/17
 * PhpStorm
 */


namespace CacheSystem\Cache\Singleton;

use Illuminate\Support\Facades\Redis;

/**
 * Class RedisSingleton
 * @package App\Components\CacheSystem\Query\Singleton
 */
final class RedisSingleton
{
    /**
     * @var \Predis\ClientInterface|null
     */
    private static $db;


    /**
     * RedisSingleton constructor.
     */
    private function __construct()
    {
    }

    /**
     *
     * @return \Predis\ClientInterface|null
     */
    public static function instance()
    {
        if (self::$db === null)

            self::$db = Redis::connection(env('REDIS_CACHE_NAME', 'chart'));

        return self::$db;
    }


}
