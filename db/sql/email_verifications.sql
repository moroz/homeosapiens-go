-- name: GetUnverifiedUserByEmail :one
select * from users where email_hash = $1 and email_confirmed_at is null;

-- name: GetExistingEmailVerificationToken :one
select * from user_tokens ut
where ut.user_id = $1 and ut.context = 'email_verification'
and ut.inserted_at > (now() - @rate_limit_period::interval);

-- name: CheckUserEmailVerificationRateLimit :one
select
    coalesce(max(ut.inserted_at) + @rate_limit_period::interval <= now(), true)::bool as can_request,
    coalesce(
        case
            when max(ut.inserted_at) + @rate_limit_period::interval > now()
            then max(ut.inserted_at) + @rate_limit_period::interval
        end,
        '0001-01-01 00:00:00'::timestamp
    )::timestamp as limited_until
from user_tokens ut
join users u on ut.user_id = u.id
where u.email_hash = $1 and u.email_confirmed_at is null
and ut.context = 'email_verification';