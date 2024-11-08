<script setup lang="ts">

</script>

<template>
    <div class="im_body">
        <router-view></router-view>
    </div>
</template>

<style lang="scss">
.im_body{
    width: 100%;
    height: 100vh;
    overflow-y: hidden;  /* 禁用滚动条 */ 
    background-color: #f0eeee;
    display: flex;
    align-items: center;
    justify-content: center;
}
</style>
