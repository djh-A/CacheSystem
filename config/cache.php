<?php
/**
 * cache.php
 * @author   djh <dengjinghui@shiyue.com>
 * @date     2021/9/17
 * PhpStorm
 * @desc:
 */


return [

    'redis' => [

        env('REDIS_CACHE_NAME', 'chart') => [
            'host' => env('REDIS_CHART_HOST', 'localhost'),
            'password' => env('REDIS_CHART_PASSWORD', null),
            'port' => env('REDIS_CHART_PORT', 6379),
            'database' => 3,
        ],
    ]

];
