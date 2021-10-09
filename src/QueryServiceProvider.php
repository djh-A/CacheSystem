<?php
/**
 * QueryServiceProvider.php
 * @author      djh <dengjinghui@shiyue.com>
 * @date     2021/9/7 20:35
 * PhpStorm
 */


namespace CacheSystem;


use CacheSystem\Query\Factory\Query;
use Illuminate\Support\ServiceProvider;

/**
 * Class QueryServiceProvider
 * @package CacheSystem
 */
class QueryServiceProvider extends ServiceProvider
{
    /**
     *
     */
    public function register()
    {
        $this->app->singleton('query', fn() => new Query());
    }

    public function boot()
    {
        $this->app->configure('cachesys');

        $path = realpath(__DIR__ . '/../config/cache.php');

        $this->mergeConfigFrom($path, 'cachesys');

    }
}
