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
import FooterContent from '~/components/FooterContent.vue';
export default {

  components: {
    HeaderNavbar,
    FooterContent
  },
  data() {
    return {
    }
  },
}
</script>

<style>
.main {
  padding: 140px 30px 0 !important;
}
</style>

