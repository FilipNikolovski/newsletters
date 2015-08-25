<?php

use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddOnDeleteCascadeInSubscribersListTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::table('subscribers_lists', function (Blueprint $table) {
            $table->dropForeign('subscribers_lists_list_id_foreign');
            $table->foreign('list_id')->references('id')->on('lists')->onDelete('cascade');
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::table('subscribers_lists', function (Blueprint $table) {
            $table->dropForeign('subscribers_lists_list_id_foreign');
            $table->foreign('list_id')->references('id')->on('lists');
        });
    }
}
