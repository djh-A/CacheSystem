<?php
/**
 * Query.php
 * @author   djh <dengjinghui@shiyue.com>
 * @date     2021/9/7
 * PhpStorm
 * @desc: Query
 */

namespace CacheSystem\Query\Factory;

use CacheSystem\Query\Producer\ImpalaSingleQuery;

class Query extends QueryFactory
{

    /**
     * @return ImpalaSingleQuery
     */
    public function ImpalaSingleQuery()
    {
        if (null === $this->impalaSingleQuery)
            $this->impalaSingleQuery = ImpalaSingleQuery::make();

        return $this->impalaSingleQuery;
    }

}
