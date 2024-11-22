<script setup lang="ts">
import {useStore} from "@/stores";


const store = useStore();
store.getUserConf()
</script>

<template>
    <div class="im_info_body">
        <div class="head">
            <router-link :to="{name: 'my_info'}">个人资料</router-link>
            <router-link :to="{name: 'base_info'}">基本信息</router-link>
            <router-link :to="{name: 'safe_info'}">安全设置</router-link>
            <router-link :to="{name: 'role_info'}">权限设置</router-link>
        </div>
        <div class="body">
            <router-view></router-view>
        </div>
    </div>

</template>

<style lang="scss">
.im_info_body{
    width: 100%;
    padding: 20px;
}

.head{
    border-bottom: 1px solid #f2f2f2;
    padding-bottom: 20px;
    a {
        margin-right: 20px;
        &.router-link-exact-active {
            color: #0099ff;
        }
    }

}

.body{
    padding: 20px 0;

    .el-form-item {
        margin-bottom: 8px;
    }
}
</style>
