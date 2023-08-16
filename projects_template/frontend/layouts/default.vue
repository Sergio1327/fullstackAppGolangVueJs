<template>
  <div>
    <HeaderNavbar />
    <main>
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
