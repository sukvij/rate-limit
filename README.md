# rate-limit

using token based bucket

/request --> 
    check if user exist or not --> 
        if not --> create  --> refill and decrease token by 1
        if yes --> check token > 0
            if yes --> calculate token again and full fill request and decrease token by 1
            if not --> calculate token again and try to full fill request
                if can then ohk
                else --> rate limit exceeds --> {total token user have, when will 1 token refill}
    
/stats
    how many token a user have and 