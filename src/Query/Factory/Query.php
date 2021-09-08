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
     * @return ImpalaQuery
     */
    public function ImpalaQuery()
    {
        if (null === $this->impala)
            $this->impala = ImpalaQuery::make();

        return $this->impala;
    }


}
