<template>
  <div>
    <HeaderNavbar />
    <main class="main">
      <div class="container">
        <nuxt></nuxt>
      </div>
    </main>
    <footerContent />
  </div>
</template>

<script>
import HeaderNavbar from '~/components/HeaderNavbar.vue';
import footerContent from '~/components/footerContent.vue';
export default {
  name: 'DefaultLayout',
  components: {
    HeaderNavbar,
    footerContent
  },
  data() {
    return {
    }
  },
}
</script>

<style>
.main {
  padding: 140px 0 0;
}
</style>

