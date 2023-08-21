<template>
    <div>
        <productModalComp class="is-flex is-align-items-center is-justify-content-flex-end" />
        <productForm />
    </div>
</template>

<script>

import productForm from '~/components/productForm.vue';
import productModalComp from '~/components/productModalComp.vue';
export default {
    components: {
        productForm,
        productModalComp
    }
}
</script>