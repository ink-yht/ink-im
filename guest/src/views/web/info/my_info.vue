<script lang="ts" setup>
import {useStore} from "@/stores";
import {nextTick, ref} from "vue";

const store = useStore();

// 修改昵称
const editNicknameRef = ref();
const showEdit = ref(false)
const editNickname = () =>{
    showEdit.value = true
    nextTick(()=>{
        editNicknameRef.value.focus()
    })
}
const editNicknameBlur = () =>{
    showEdit.value = false
    // 调更新接口
}

</script>
<template>
    <div class="my_info_view">
        <el-form-item label="头像">
            <el-avatar :src="store.userConfInfo.avatar"></el-avatar>
        </el-form-item>
        <el-form-item label="用户号">
            <span>{{store.userConfInfo.id}}</span>
        </el-form-item>
        <el-form-item label="昵称">
            <span v-if="!showEdit">{{store.userConfInfo.nickname}}</span>
            <el-input ref="editNicknameRef" :maxlength="16" v-else class="edit_nickname_ipt" v-model="store.userConfInfo.nickname" placeholder="修改昵称"></el-input>
            <i class="iconfont icon-edit" @click="editNickname"></i>
        </el-form-item>
        <el-form-item label="简介">
            <span>{{store.userConfInfo.abstract}}</span>
            <i class="iconfont icon-edit"></i>
        </el-form-item>
    </div>
</template>
<style >
.my_info_view{
    i {
        cursor: pointer;
        margin-left: 7px;
    }

    .edit_nickname_ipt{
        display: inline-flex;
        width: 200px;
    }
}
</style>