<template>
    <div>
        <button @click="openModal" value="" type="submit">Добавить продажу</button>
        <SaleModal v-if="modalVisible" @dataSubmitted="submitModalData" @closeModal="closeModal" />
    </div>
</template>

<script>
import SaleModal from './SaleModal.vue';
export default {
    components: {
        SaleModal
    },
    methods: {
        openModal() {
            this.modalVisible = true;
        },
        submitModalData() {
            // Здесь можно обновить список данных
            this.modalVisible = false;
        },
        closeModal() {
            this.modalVisible = false;
        },
    },
    data() {
        return {
            saleListData: [],
            modalVisible: false,
        }
    }
}
</script>