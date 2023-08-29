<template>
    <div class="is-flex is-align-items-center is-justify-content-flex-end">
        <button class="btn" @click="openModal" value="" type="submit">Добавить продажу</button>
        <addSaleForm v-if="modalVisible" @closeModal="closeModal" :modalVisible="modalVisible" />
    </div>
</template>

<script>
import addSaleForm from './addSaleForm.vue';
export default {
    components: {
        addSaleForm
    },
    methods: {
        openModal() {
            this.modalVisible = true;
        },
        submitModalData() {
            this.modalVisible = false;
        },
        closeModal() {
            this.modalVisible = false;
        },
    },
    data() {
        return {
            modalVisible: false,
        }
    }
}
</script>

<style>
.btn {
    padding: 14px 28px;
    background-color: rgb(255, 190, 70) !important;
    border-radius: 5px;
}
</style>