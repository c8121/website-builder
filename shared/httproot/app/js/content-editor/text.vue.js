export default {
    template: `
        <div>
            <textarea v-model="editorContent" class="form-control"></textarea>
        </div>
    `,

    data() {
        return {
            editorContent: null
        }
    },


    methods: {
        setContent(content) {
            this.editorContent = content;
        },

        getContent() {
            return this.editorContent;
        },
    },

    emits: ['componentEvent'],
}