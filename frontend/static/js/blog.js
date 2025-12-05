class Blog{
    constructor() {
        this.currentUser = null;
        this.currentPage = 1;
        this.commentsPerPage = 10;
        this.init();
    }

    init() {
        this.bindEvents();
        //this.checkLoginStatus();
        //this.loadComments();
    }

    async loadBlogPost(postID) {
        try {
            // 获取文章元数据
            const metaResponse = await fetch(`/api/v0/post/${postID}`);
            const articleMeta = await metaResponse.json();

            console.log('Article Metadata:', articleMeta);
            console.log('Author Name:', articleMeta.AuthorName);
            console.log('Last Update:', articleMeta.LastUpdated);

            // Update article metadata in the DOM
            document.getElementById('article-title').textContent = articleMeta.Title || 'Unknown Title';
            document.getElementById('article-author').textContent = articleMeta.AuthorName || 'Unknown Author';
            document.getElementById('article-date').textContent = articleMeta.LastUpdated || 'Unknown Date';
            
            // Fetch Markdown content
            const mdResponse = await fetch(articleMeta.ContentFilePath);
            console.log('Markdown Content Path:', articleMeta.ContentFilePath);
            console.log('Markdown Response:', mdResponse);
            const markdownContent = await mdResponse.text();
            
            // Render Markdown to HTML
            document.getElementById('markdown-content').innerHTML = DOMPurify.sanitize(await marked.parse(markdownContent));
            
            
        
            
        } catch (error) {
            console.error('Failed to load blog post:', error);
            document.getElementById('markdown-content').innerHTML = '<p>Failed to load markdown content</p>';
        }
    }

    bindEvents() {
    }
}

export { Blog };