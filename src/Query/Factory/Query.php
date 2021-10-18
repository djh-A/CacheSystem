<?php
/**
 * Query.php
 * @author   djh <dengjinghui@shiyue.com>
 * @date     2021/9/7
 * PhpStorm
 * @desc: Query
 */

namespace CacheSystem\Query\Factory;

use CacheSystem\Query\Producer\ImpalaQuery;


/**
 * Class Query
 * @package CacheSystem\Query\Factory
 */
class Query extends QueryFactory
{

    /**
     * Impala query
     *
     * @param null $sql
     * @return ImpalaQuery|mixed
     */
    public function ImpalaQuery($sql = null)
    {
        if (null === $this->impala)
            $this->impala = ImpalaQuery::make();

        if ($sql)
            return $this->impala->get($sql);

        return $this->impala;
    }


}
