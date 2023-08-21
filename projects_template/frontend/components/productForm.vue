<template>
    <div>
        <div class="field">
            <label class="label">Название продукта или тег продукта</label>
            <div class="control">
                <input class="input" type="text" v-model="productName">
            </div>
        </div>
        <div class="field">
            <label class="label">Лимит вывода</label>
            <div class="control">
                <div class="select">
                    <select v-model="limit">
                        <option value="1">1</option>
                        <option value="5">5</option>
                        <option value="10">10</option>
                        <option value="50">50</option>
                    </select>
                </div>
            </div>
        </div>

    </div>
</template>
<script>
export default {
}
</script>