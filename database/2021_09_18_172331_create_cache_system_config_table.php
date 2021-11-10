<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

class CreateCacheSystemConfigTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('cache_system_configs', function (Blueprint $table) {
            $table->bigIncrements('id');
            $table->text('sql')->comment('sql');
            $table->string('cacheKey')->comment('缓存的key');
            $table->tinyInteger('isCache')->default(1)->comment('1 缓存 0 不缓存');
            $table->integer('expTime')->default(300)->comment('缓存的有效期,默认5分钟,秒为单位');
            $table->float('data_size')->default(0.00)->comment('缓存的数据大小(KB)');
            $table->float('query_time_consuming')->default(0.00)->comment('该条sql查询数仓时的耗时');
            $table->float('cache_time_consuming', 10, 6)->default(0.0000)->comment('该条sql查询缓存时的耗时');
            $table->bigInteger('heat')->default(1)->comment('缓存的热度值');
            $table->dateTime('effective_date')->comment('有效日期');
            $table->string('route')->comment('路由');
            $table->integer('product_id')->default(0)->comment('游戏id');
            $table->timestamps();
            $table->index('isCache');
            $table->index(['product_id', 'route']);
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('cache_system_configs');
    }
}
