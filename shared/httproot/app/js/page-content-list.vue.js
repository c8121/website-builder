import { nextTick } from 'vue'

export default {
    template: `
        <div>
            <table class="table">
                <tbody>
                    <tr v-for="(item, index) in list" @click="itemClicked(item)" role="button">
                        <td>{{ item.Name }}</td>
                        <td>{{ item.ContentType }}</td>
                        <td class="text-end text-nowrap">
                            <button @click.stop="move(index, -1)" v-if="index > 0" class="btn btn-sm m-0 p-0" title="Move up"><i class="bi bi-arrow-up-circle"></i></button>
                            <button @click.stop="move(index, 1)" v-if="index < list.length - 1" class="btn btn-sm m-0 p-0 ms-1" title="Move down"><i class="bi bi-arrow-down-circle"></i></button>
                        </td>
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

        move(index, delta) {
            const self = this;

            const target = index + delta;
            if (target < 0 || target >= self.list.length)
                return;

            const ids = self.list.map(item => item.Id);
            [ids[index], ids[target]] = [ids[target], ids[index]];

            const requestParams = {
                method: 'POST',
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify({PageId: self.value ? self.value : 0, Ids: ids})
            }
            fetch('/content/reorder', requestParams)
                .finally(() => self.loadItems());
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