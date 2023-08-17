<script>
import formVue from "~/components/SaleForm.vue";
import saleList from "~/components/saleList.vue";
export default {
    components: {
        formVue,
        saleList
    }, methods: {
        handleData(data) {
            this.saleListData = data
        }
    }, data() {
        return {
            saleListData: []
        }
    }
}
</script>

<template>
    <div>
        <formVue @data="handleData" />
        <saleList :saleList="saleListData" />
    </div>
</template>