<template>
    <header class="navbar is-dark py-2 px-1">
        <div class="container">
            <div class="navbar-brand">
                <nuxt-link to="/">
                    <img src="../assets/logo.png" class="image is-64x64 is-rounded img" alt="">
                </nuxt-link>

                <a role="button" class="navbar-burger burger" aria-label="menu" aria-expanded="false"
                    data-target="navbarMenu" @click="toggleNavbar">
                    <span aria-hidden="true"></span>
                    <span aria-hidden="true"></span>
                    <span aria-hidden="true"></span>
                </a>
            </div>
            <div class="navbar-menu" :class="{ 'is-active': isNavbarActive }" id="navbarMenu">
                <div class="navbar-end">
                    <nuxt-link class="navbar-item " to="/" :class="{ 'is-active': isActive('/') }">
                        Главная
                    </nuxt-link>
                    <nuxt-link class="navbar-item" to="/productList" :class="{ 'is-active': isActive('/productList') }">
                        Продукты
                    </nuxt-link>
                    <nuxt-link class="navbar-item" to="/stockList" :class="{ 'is-active': isActive('/stockList') }">
                        Склады
                    </nuxt-link>
                    <nuxt-link class="navbar-item" to="/saleList" :class="{ 'is-active': isActive('/saleList') }">
                        Продажи
                    </nuxt-link>
                </div>
            </div>
        </div>
    </header>
</template>
  
<script>
export default {
    data() {
        return {
            isNavbarActive: false
        };
    },
    methods: {
        toggleNavbar() {
            this.isNavbarActive = !this.isNavbarActive;
        },
        isActive(route) {
            return this.$route.path === route;
        }
    }
};
</script>
  
<style scoped>
.img {
    border-radius: 50%;
}

.navbar {
    position: fixed !important;
    width: 100%;
}

.navbar-item.is-active {
    background-color: orange !important;
}
</style>