<?php
/**
 * QueryInterface.php
 * @author   djh <dengjinghui@shiyue.com>
 * @date     2021/9/7
 * PhpStorm
 * @desc:
 */

namespace CacheSystem\Query\Producer;

interface QueryInterface
{
    /**
     * @param string|array $sql
     * @return mixed
     */
    public function get($sql);

}
