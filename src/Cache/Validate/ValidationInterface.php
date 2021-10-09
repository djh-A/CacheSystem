<?php
/*
 * Copyright (c) 2021. Ownership belongs to Deng JingHui.
 * The date and time when the current file was last changed is 2021/4/27 上午10:55.
 * The name of the currently opened file where the notice is to be generated is ValidationInterface.php.
 * ValidationInterface.php
 * @author   djh <dengjinghui@shiyue.com>
 * @date     2021/4/27
 * PhpStorm
 * @desc: Rule verification interface
 */


namespace CacheSystem\Cache\Validate;


use CacheSystem\Cache\Producer\SysCache;

/**
 * Interface ValidationInterface
 * @package App\Components\CacheSystem\Cache\Validate
 */
interface ValidationInterface
{
    /**
     * @param SyCache $cache
     * @param $result
     * @return mixed
     * @throws ValidationException
     */
    public function rule(SysCache $cache, $result);
}
