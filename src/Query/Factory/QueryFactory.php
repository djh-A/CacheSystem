<?php
/**
 * QueryFactory.php
 * @author   djh <dengjinghui@shiyue.com>
 * @date     2021/9/7
 * PhpStorm
 * @desc: Query abstract class
 */

namespace CacheSystem\Query\Factory;

use CacheSystem\Query\Producer\QueryInterface;

/**
 * Class QueryFactory
 * @package CacheSystem\Query\Factory
 */
abstract class QueryFactory
{

    /**
     * @var QueryInterface
     */
    protected $impala;

    /**
     * @param null $sql
     * @return mixed
     */
    abstract public function ImpalaQuery($sql = null);

}
