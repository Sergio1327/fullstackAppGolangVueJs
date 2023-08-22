<template>
    <div class="modal is-active">
        <div class="modal-background"></div>
        <div class="modal-content box">
            <div class="field">
                <label class="label">Название продукта</label>
                <div class="control">
                    <input class="input" v-model="productName" type="text">
                </div>
            </div>
            <div class="field">
                <label class="label">Описание продукта</label>
                <div class="control">
                    <input class="input" v-model="descr" type="text">
                </div>
            </div>
            <div class="field">
                <label class="label">Варианты продукта</label>
                <div class="control">
                    <div v-for="(variant, index) in productVariants" :key="index">
                        <div class="field">
                            <label class="label">Масса</label>
                            <div class="control">
                                <input class="input" v-model="variant.weight" type="text">
                            </div>
                        </div>
                        <div class="field">
                            <label class="label">Единица измерения</label>
                            <div class="control">
                                <input class="input" v-model="variant.unit" type="text">
                            </div>
                        </div>
                    </div>
                    <button class="button is-primary" type="button" @click="addVariant">Добавить вариант</button>
                </div>
            </div>
            <div class="label">{{ resp }}</div>
            <button class="button is-primary" @click="submitModalData">Добавить продажу</button>
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
            descr: "",
            productVariants: [{
                weight: "",
                unit: ""
            }],
            resp: ""
        };
    },
    methods: {
        closeModal() {
            this.$emit('closeModal');
        },
        addVariant() {
            this.productVariants.push({});
        },
        submitModalData() {
            // Здесь можно выполнить логику отправки данных на сервер,
            // используя данные this.productName, this.descr и this.productVariants
            console.log('Отправка данных на сервер:', {
                productName: this.productName,
                descr: this.descr,
                productVariants: this.productVariants
            });
        }
    }
}
</script>
  