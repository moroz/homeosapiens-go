package main

const listVideosQuery = `
select locale, id, title,
    'Dr ' || host host,
	case host
		when 'Sanjay Modi' then '/modi.jpeg'
		when 'Asher Shaikh' then '/dr-asher.jpeg'
		when 'Herman Jeggels' then '/jeggels.jpeg'
	end profile_picture,
recorded_on::text from (
	select 'en' locale, v.id, v.title_en title, h.given_name || ' ' || h.family_name host, recorded_on from videos v
	join hosts h on v.host_id = h.id
	union all
	select 'pl' locale, v.id, v.title_pl, h.given_name || ' ' || h.family_name, recorded_on from videos v
	join hosts h on v.host_id = h.id
) s order by s.id;
`

const insertAssetQuery = `
insert into assets (id, object_key) values ($1, $2)
`

const setThumbnailQuery = `
update videos
set
thumbnail_en_id = case $1 when 'pl' then thumbnail_en_id else $2 end,
thumbnail_pl_id = case $1 when 'pl' then $2 else thumbnail_pl_id end
where id = $3;
`
