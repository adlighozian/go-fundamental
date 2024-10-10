package repository

const (
	getPosts = `select 
		p.id, 
		p.title,
		p.content, 
		p.status, 
		p.publish_date  
	from posts p`

	getTags = `select 
		t.id , 
		t."label"  
	from posts p 
	join post_tags pt on pt.post_id = p.id 
	join tags t on t.id = pt.tag_id 
	where 1=1
		and p.id = $1; `

	getIdPostsTags = `select 
		pt.id  
	from posts p 
	join post_tags pt on pt.post_id = p.id 
	join tags t on t.id = pt.tag_id 
	where 1=1
		and p.id = $1; `

	detailPosts = `select 
		p.id, 
		p.title,
		p.content, 
		p.status, 
		p.publish_date  
	from posts p
	where id = $1`

	checkDetailPosts = `select 
		count(*) 
	from posts p
	where id = $1`

	deletePosts = `delete from posts where id = $1 `

	deletePostsTags = `delete from post_tags where id = $1`

	createPost = `INSERT 
	INTO posts (title, content, status)
	VALUES ($1, $2, 'draft')
	RETURNING id;`

	getIdTag = `select
		id
	from tags
	where LOWER(label) = LOWER($1)`

	creatPostsTag = `insert 
	into post_tags (post_id , tag_id) 
	VALUES ($1, $2)`

	updatePosts = `UPDATE posts
	SET title = $1, content = $2, status = 'publish', publish_date = NOW(), updated_at = NOW()
	WHERE id = $3`

	loginSql = `SELECT 
		username, 
		password, 
		role 
	FROM users WHERE username=$1`

	registerSql = `INSERT 
	INTO users (username, password, role) 
	VALUES ($1, $2, $3)`
)
