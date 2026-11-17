export default {
    template: `
        <div>
            <textarea style="height: 50vh" v-model="editorContent" class="form-control"></textarea>
        </div>
    `,

    data() {
        return {
            editorContent: null
        }
    },

    methods: {
        setValue(value) {
            this.editorContent = value
        },

        getValue() {
            return this.editorContent;
        },
    },
}