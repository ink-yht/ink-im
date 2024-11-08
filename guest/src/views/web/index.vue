<script setup lang="ts">
import { Menu, ChatLineRound, User, Notification } from "@element-plus/icons-vue"
import { useRouter,useRoute } from "vue-router";

const router = useRouter()
const route = useRoute()

const checkMenu = (menu: string) => {
    router.push({
        name: menu,
    })
};


</script>

<template>
    <div class="im_web">
        <div class="im_slide">
            <div class="avatar" @click="checkMenu('info')">
                <img src="https://img2.baidu.com/it/u=186808850,2178610585&fm=253&app=120&size=w931&n=0&f=JPEG&fmt=auto?sec=1731171600&t=99cdcbcb995f0802732e92b5a8a0eb58"
                    alt="">
            </div>
            <div class="im_menus">
                <div class="icon" :class="{active: route.name === 'session'}" @click="checkMenu('session')">
                    <el-icon>
                        <ChatLineRound />
                    </el-icon>
                </div>
                <div class="icon" :class="{active: route.name === 'contacts'}" @click="checkMenu('contacts')">
                    <el-icon>
                        <User />
                    </el-icon>
                </div>
                <div class="icon" :class="{active: route.name === 'notice'}" @click="checkMenu('notice')">
                    <el-icon>
                        <Notification />
                    </el-icon>
                </div>
            </div>
            <div class="other icon" :class="{active: route.name === 'info'}" @click="checkMenu('info')">
                <el-icon>
                    <Menu />
                </el-icon>
            </div>
        </div>
        <router-view></router-view>
    </div>

</template>

<style lang="scss">
.im_web {
    width: 1000px;
    height: 600px;
    background-color: white;
    border-radius: 10px;
    overflow: hidden;
    box-shadow: 0 0 5px 3px rgba(0, 0, 0, 0.05);
    display: flex;
}

.im_slide {
    width: 80px;
    height: 100%;
    background-color: #0099ff;
    position: relative;
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 20px 0;

    .avatar {
        width: 46px;
        height: 46px;

        img {
            width: 100%;
            height: 100%;
            border-radius: 50%;
        }
    }

    .im_menus {
        margin-top: 20px;
        display: flex;
        flex-direction: column;
        align-items: center;

        >div {
            margin-bottom: 20px;
        }

    }

    .icon {
        cursor: pointer;
        font-size: 22px;
        color: white;
        width: 40px;
        height: 40px;
        border-radius: 5px;
        display: flex;
        justify-content: center;
        align-items: center;
        transition: all 0.3s;

        &:hover {
            background-color: #33adff;
        }

        &.active {
            background-color: #33adff;
        }
    }

    .other {
        position: absolute;
        bottom: 20px;
    }
}
</style>