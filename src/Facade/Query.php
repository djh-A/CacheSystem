<?php
/**
 * Query.php
 * @author      djh <dengjinghui@shiyue.com>
 * @date     2021/9/7 20:29
 * PhpStorm
 */


namespace CacheSystem\Facade;


use CacheSystem\Query\Producer\QueryInterface;
use Illuminate\Support\Facades\Facade;

/**
 * Class Query
 * @package CacheSystem\Facade
 * @method static QueryInterface ImpalaQuery()
 */
class Query extends Facade
{
    /**
     * @return string
     */
    public static function getFacadeAccessor()
    {
        return 'query';
    }
}
