<?php
/*
 * Copyright (c) 2021. Ownership belongs to Deng JingHui.
 * The date and time when the current file was last changed is 2021/4/27 上午10:46.
 * The name of the currently opened file where the notice is to be generated is ValidationException.php.
 * ValidationException.php
 * @author      djh <dengjinghui@shiyue.com>
 * @date     2021/4/27 10:46
 * PhpStorm
 */


namespace CacheSystem\Cache\Validate;


/**
 * Class ValidationException
 * @package App\Components\CacheSystem\Cache\Validate
 */
class ValidationException extends \Exception
{


    /**
     * ValidationException constructor.
     */
    public function __construct($message)
    {
        parent::__construct($message);
    }
}
