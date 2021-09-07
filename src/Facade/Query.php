<?php
/**
 * Query.php
 * @author      djh <dengjinghui@shiyue.com>
 * @date     2021/9/7 20:29
 * PhpStorm
 */


namespace CacheSystem\Facade;


use Illuminate\Support\Facades\Facade;

/**
 * Class Query
 * @package CacheSystem\Facade
 */
class Query extends Facade
{
    /**
     * @return string
     */
    public function getFacadeAccessor()
    {
        return 'query';
    }
}
