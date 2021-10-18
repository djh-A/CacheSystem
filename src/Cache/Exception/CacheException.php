<?php
/**
 * CacheException.php
 * @author      djh <dengjinghui@shiyue.com>
 * @date     2021/9/17 17:53
 * PhpStorm
 */


namespace CacheSystem\Cache\Exception;


/**
 * Class CacheException
 * @package CacheSystem\Cache\Exception
 */
class CacheException extends \Exception
{

    /**
     * CacheException constructor.
     */
    /**
     * CacheException constructor.
     * @param string $message
     * @param int $code
     * @param Throwable|null $previous
     */
    public function __construct($message = "", $code = 0, Throwable $previous = null)
    {
        parent::__construct($message, $code, $previous);
    }
}
