<template>
    <div class="card mb-5  py-3 ">
        <div class="card-content is-flex is-align-items-center">
            <div class="media m-0">
                <div class="media-left">
                    <figure class="image is-96x96">
                        <img :src="card.imgLink" alt="Placeholder image">
                    </figure>
                </div>

            </div>
            <div class="media-content">
                <p class="title is-4 m-0">{{ card.title }}</p>
                <div class="content">
                    {{ card.descr }}
                </div>
            </div>

        </div>
    </div>
</template>

<script>

export default {
    props: {
        card: {
            type: Object,
            required: true
        }
    }
}
</script>