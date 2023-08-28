<script>
import SaleCardVue from './SaleCard.vue'

export default {
    components: {
        SaleCardVue
    },
    props: {
        saleList: {
            type: Array,
            required: true,
        }
    }
}
</script>

<template>
    <div class="mt-6 ">
        <SaleCardVue class="mb-6" v-for="s in saleList" :key="s.SaleID" :data="s" />
    </div>
</template> 