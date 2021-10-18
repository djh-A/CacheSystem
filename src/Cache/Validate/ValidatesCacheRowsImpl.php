<?php
/*
 * Copyright (c) 2021. Ownership belongs to Deng JingHui.
 * The date and time when the current file was last changed is 2021/4/27 上午11:14.
 * The name of the currently opened file where the notice is to be generated is ValidatesCacheRowsImpl.php.
 * ValidatesCacheRowsImpl.php
 * @author      djh <dengjinghui@shiyue.com>
 * @date     2021/4/27 11:13
 * PhpStorm
 */


namespace CacheSystem\Cache\Validate;


use CacheSystem\Cache\Producer\SysCache;

/**
 * Class ValidatesCacheRowsImpl
 * @package App\Components\CacheSystem\Cache\Validate
 */
class ValidatesCacheRowsImpl implements ValidationInterface
{
    /**
     * @param SyCache $cache
     * @param $result
     * @throws ValidationException
     */
    public function rule(SysCache $cache, $result)
    {
        if (count((array)$result) > $cache->getMaxCacheDataRows()) {
            throw new ValidationException(sprintf("The number of rows in the result set exceeds the limit,
            the maximum number of rows currently allowed is:%s", $cache->getMaxCacheDataRows()));
        }
    }
}
