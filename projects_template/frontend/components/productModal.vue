<template>
    <div class="modal is-active">
        <div class="modal-background"></div>
        <div class="modal-content box">
            <div class="field">
                <label class="label">Название продукта</label>
                <div class="control">
                    <input class="input" v-model="productName" type="tex">
                </div>
            </div>
            <div class="field">
                <label class="label">ID склада</label>
                <div class="control">
                    <input class="input" type="number">
                </div>
            </div>
            <div class="field">
                <label class="label">Кол-во продуктов</label>
                <div class="control">
                    <input class="input" type="number">
                </div>
            </div>
            <div class="label"> {{ resp }}</div>
            <button class="button is-primary" @click="submitModalData">добавить продажу</button>
            <button class="button" @click="closeModal">Закрыть</button>
        </div>
        <button class="modal-close is-large" aria-label="Закрыть" @click="closeModal"></button>

    </div>
</template>

<script>
export default {
    data() {
        return {
            productName: "",
            resp: ""
        };
    },
    methods: {
        closeModal() {
            this.$emit('closeModal');
        },
    }
}
</script>
