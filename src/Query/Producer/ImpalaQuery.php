<?php
/**
 * ImpalaQuery.php
 * @author      djh <dengjinghui@shiyue.com>
 * @date     2021/9/8 14:11
 * PhpStorm
 */


namespace CacheSystem\Query\Producer;


use App\Components\DwQueryClient;
use App\Components\ThriftQuery\ThriftQuery;
use App\Models\Remote\Product;

/**
 * Class ImpalaQuery
 * @package CacheSystem\Query\Producer
 */
class ImpalaQuery extends BaseQuery implements QueryInterface
{

    /**
     * @inheritDoc
     */
    public function get($sql)
    {

        if (is_callable($sql))
            $sql = call_user_func($sql);

        if (is_array($sql))
            return DwQueryClient::query($sql);

        return (new ThriftQuery())->queryAll($sql);
    }
}
