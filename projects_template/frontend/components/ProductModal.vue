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
                <label class="label">Дата добавления продукта</label>
                <div class="control">
                    <input class="input" type="datetime-local" v-model="addedAt">
                </div>
            </div>
            <div class="field">
                <label class="label">теги продукта</label>
                <div class="control">
                    <input class="input" v-model="tags" type="text">
                </div>
            </div>
            <div class="field">
                <label class="label">Варианты продукта</label>
                <div class="control">
                    <div v-for="(variant, index) in productVariants" :key="index">
                        <div class="field">
                            <label class="label">Масса</label>
                            <div class="control">
                                <input class="input" v-model="variant.weight" type="number">
                            </div>
                        </div>
                        <div class="field">
                            <label class="label">Единица измерения</label>
                            <div class="control">
                                <input class="input" v-model="variant.unit" type="text">
                            </div>
                        </div>
                    </div>
                    <div class="field is-flex is-justify-content-flex-end is-align-items-center">
                        <button class="button is-primary mt-4" type="button" @click="addVariant">Добавить вариант</button>
                    </div>
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
            addedAt: null,
            tags: "",
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
        async submitModalData() {
            const formatedVariants = this.productVariants.map(v => ({
                weight: parseInt(v.weight),
                unit: v.unit
            }))

            const requestData = {
                name: this.productName,
                description: this.descr,
                added_at: this.formatDate(this.addedAt),
                tags: this.tags,
                variants: formatedVariants
            }
            console.log(JSON.stringify(requestData))
            try {
                const response = await fetch("http://127.0.0.1:9000/product/add", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(requestData),
                })

                const responseData = await response.json()
                console.log(responseData.Data)
                this.resp = "продукт успешно добавлен в базу , ID продукта = " + responseData.Data.product_id
            } catch (error) {
                console.error("Ошибка при отправке запроса:", error);
                console.log(error)
            }
        },
        formatDate(dateTime) {
            const date = new Date(dateTime);
            const year = date.getUTCFullYear();
            const month = String(date.getUTCMonth() + 1).padStart(2, "0");
            const day = String(date.getUTCDate()).padStart(2, "0");
            const hours = String(date.getHours()).padStart(2, "0");
            const minutes = String(date.getUTCMinutes()).padStart(2, "0");
            const seconds = String(date.getUTCSeconds()).padStart(2, "0");

            return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}+05:00`;
        }
    }
}
</script>
  