<?php
/**
 * ImpalaSingleQuery.php
 * @author      djh <dengjinghui@shiyue.com>
 * @date     2021/9/7 20:10
 * PhpStorm
 * @desc impala Single sql query
 */


namespace CacheSystem\Query\Producer;


use App\Components\ThriftQuery\ThriftQuery;

/**
 * Class ImpalaSingleQuery
 * @package CacheSystem\Query\Producer
 */
class ImpalaSingleQuery extends BaseQuery implements QueryInterface
{

    /**
     * @inheritDoc
     */
    public function get($sql)
    {
        if (is_callable($sql))
            $sql = call_user_func($sql);

        return (new ThriftQuery())->queryAll($sql);
    }
}
