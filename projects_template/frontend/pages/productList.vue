<template>
    <div>
        <productModalComp class="is-flex is-align-items-center is-justify-content-flex-end" />
        <productForm @productData="handleData" />
        <productCardList :productList="productList" />
    </div>
</template>

<script>

import productForm from '~/components/ProductForm.vue';
import productModalComp from '~/components/ProductModalComp.vue';
import productCardList from '~/components/ProductCardList.vue';

export default {
    components: {
        productForm,
        productModalComp,
        productCardList
    },
    methods: {
        handleData(productData) {
            this.productList = productData
        }
    },
    data() {
        return {
            productList: []
        }
    }
}
</script>