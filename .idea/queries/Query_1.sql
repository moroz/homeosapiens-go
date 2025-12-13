select v.* from videos v
where v.is_public = true
or v.id in (select v.id from videos v
                        join event_registrations er on v.event_id = er.event_id
                        join users u on er.user_id = u.id
                        where u.email = 'karol@moroz.dev')
or exists (select 1 from users u where u.email = 'karol@moroz.dev' and u.user_role = 'Administrator')