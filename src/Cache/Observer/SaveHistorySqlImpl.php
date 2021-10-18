<?php
/**
 * SaveHistorySqlImpl.php
 * @author      djh <dengjinghui@shiyue.com>
 * @date     2021/9/22 14:55
 * PhpStorm
 */


namespace CacheSystem\Cache\Observer;


use App\Jobs\ExampleJob;
use App\Jobs\SaveHistorySqlJob;
use App\Models\CacheSystemConfig;
use CacheSystem\Cache\Producer\SysCache;
use Illuminate\Support\Carbon;
use Illuminate\Support\Facades\Request;

class SaveHistorySqlImpl implements CacheObserver
{
    protected $isQueue = true;

    /**
     * @inheritDoc
     */
    public function onChange(...$args)
    {
        /**
         * @var SysCache $cache
         */
        $cache = $args[0][0];

        $cache->setRoute(Request::path());

        $cacheKey = (array)$cache->getCacheDriver()->getCacheKey();

        if ($this->isQueue) {

            dispatch((new SaveHistorySqlJob($cache))->onQueue('default'));
        } else {

            foreach ($cacheKey as $key => $cKey) {

                if ($query = CacheSystemConfig::query()->where('cacheKey', $cKey)->first()) {

                    $query->increment('heat');

                    $attr = [
                        'data_size' => $cache->getCacheDriver()->getDataSize()[$key] ?? 0,
                        'cache_time_consuming' => $cache->getCacheDriver()->getUseTime(),
                        'route' => $cache->getRoute(),
                        'effective_date' => Carbon::today()->toDateString()
                    ];

                    if ($useTime = $cache->getUseTime())
                        $attr = array_merge($attr, ['query_time_consuming' => $useTime]);

                    $query->update($attr);

                } else {

                    $attr = [
                        'cacheKey' => $cKey,
                        'sql' => $cache->getSql()[$key],
                        'data_size' => $cache->getCacheDriver()->getDataSize()[$key] ?? 0,
                        'cache_time_consuming' => $cache->getCacheDriver()->getUseTime(),
                        'route' => $cache->getRoute(),
                        'effective_date' => Carbon::today()->toDateString()
                    ];

                    if ($useTime = $cache->getUseTime())
                        $attr = array_merge($attr, ['query_time_consuming' => $useTime]);


                    CacheSystemConfig::create($attr);
                }

            }

        }
    }
}
