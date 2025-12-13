select v.* from users u, videos v
left join events e on v.event_id = e.id
left join event_registrations er on er.event_id = e.id
where u.email = 'karol@moroz.dev' and (v.is_public = true or u.user_role = 'Administrator')