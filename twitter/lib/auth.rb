module Auth
    def auth_test
        "Returning NOW FOR TAI COW from Auth Module, aye!"
    end

    def get_session_id(session)
        return session.id
    end
    
    def start_session(session, user_id)
        # puts params
        # return params
        session[:current_user_id] = user_id
        session[:session_active] = true
        return session.id
    end

    def end_session(session)
        session.clear
        reset_session
    end

    def is_logged_in(session, session_id, user_id)
        result = {
            "success": false,
            "message": "What went wrong?"
        }
        if session.id.to_s === session_id.to_s
            if session[:current_user_id].to_s === user_id.to_s
                result["success"] = true
                result["message"] = "Logged in!"
            end
        else
            result["success"] = false
            result["message"] = "Can't find session."
            forceLogout(session)
        end

        return result
    end

    def forceLogout(session)
        session.clear
        session[:session_active] = false
    end

end