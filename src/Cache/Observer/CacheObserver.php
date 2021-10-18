<?php
/**
 * CacheObserver.php
 * @author      djh <dengjinghui@shiyue.com>
 * @date     2021/9/18 16:16
 * PhpStorm
 */


namespace CacheSystem\Cache\Observer;

/**
 * Interface CacheObserver
 * @package CacheSystem\Cache\Observer
 */
interface CacheObserver
{
    /**
     * @param mixed ...$args
     * @return mixed
     */
    public function onChange(...$args);
}
