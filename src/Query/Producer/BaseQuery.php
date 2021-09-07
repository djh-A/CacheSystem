<?php
/**
 * BaseQuery.php
 * @author   djh <dengjinghui@shiyue.com>
 * @date     2021/9/7
 * PhpStorm
 * @desc:
 */

namespace CacheSystem\Query\Producer;

abstract class BaseQuery
{
    /**
     * BaseQuery constructor.
     */
    private function __construct()
    {
    }


    /**
     * @return mixed
     */
    public static function make()
    {
        return static();
    }
}
