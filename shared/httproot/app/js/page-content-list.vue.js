import { nextTick } from 'vue'

export default {
    template: `
        <div>
            <table class="table">
                <tbody>
                    <tr v-for="item in list" @click="itemClicked(item)" role="button">
                        <td>{{ item.Name }}</td>
                        <td>{{ item.ContentType }}</td>
                    </tr>
                </tbody>
            </table>
        </div>
    `,

    props: {
        value: "",
    },

    data() {
        return {
            pageId: 0,
            list: []
        }
    },

    watch: {
        value() {
            this.loadItems();
        }
    },

    methods: {

        loadItems() {
            const self = this;

            const requestOptions = {
                method: 'GET'
            }
            const pageId = self.value ? self.value : 0;
            fetch('/content/list/' + pageId, requestOptions)
                .then(res => res.json())
                .then(json => {
                    self.list.splice(0, self.list.length);
                    for (const item of json) {
                        self.list.push(item);
                    }
                });
        },

        reload() {
            this.loadItems();
        },

        itemClicked(item) {
            this.$emit('componentEvent', 'pageContentClick', 'page-content-list', item);
        },
    },

    emits: ['componentEvent'],

    created() {
        this.loadItems();
    }
}