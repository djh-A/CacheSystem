<?php
/*
 * Copyright (c) 2021. Ownership belongs to Deng JingHui.
 * The date and time when the current file was last changed is 2021/4/27 上午11:09.
 * The name of the currently opened file where the notice is to be generated is Config.php.
 * Config.php
 * @author      djh <dengjinghui@shiyue.com>
 * @date     2021/4/27 11:09
 * PhpStorm
 */


namespace CacheSystem\Cache\Validate;


/**
 * Class Config
 * @package App\Components\CacheSystem\Cache\Validate
 */
class Config
{
    /**
     * Verification class
     *
     * @return string[]
     */
    public static function verifier()
    {
        return [
            //Verify the number of cached result rows
            ValidatesCacheRowsImpl::class,
        ];
    }
}
